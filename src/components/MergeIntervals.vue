<template>
  <v-container>
    <h1>Merge Intervals</h1>
    <v-alert id="error" v-if="errors.length" type="error" class="mt-8">
      <ul>
        <li v-for="(error, index) in errors" :key="index">{{ error }}</li>
      </ul>
    </v-alert>
    <v-alert id="info" v-if="errors.length" type="info" class="mt-8">
      <ul>
        <li>Intervals must consist of two integers, separated by a comma.</li>
        <li>The first integer in the interval must not be greater than the second.</li>
        <li>Negative integers are allowed.</li>
        <li>Whitespace makes no difference.</li>
      </ul>
    </v-alert>
    <v-form class="pt-8" @submit.prevent="submitForm">
      <v-textarea
        label="Intervals"
        v-model="intervalsInput"
        placeholder="[25,30] [2,19] [14,23] [4,8]"
      ></v-textarea>
      <v-btn id="merge" color="green" class="mt-4" type="submit">Merge</v-btn>
      <v-btn id="reset" color="blue" class="mt-4 ml-4" @click="resetForm">Reset</v-btn>
    </v-form>
    <h3 id="result_header" v-if="mergedIntervals.length > 0" class="mt-6">Result:</h3>
    <p v-if="mergedIntervals.length > 0" class="mt-2" id="results">
      <span v-for="(interval, index) in mergedIntervals" :key="index">
        {{ `[${interval.start},${interval.end}]` }}
      </span>
    </p>
  </v-container>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import axios from 'axios';

export default defineComponent({
  name: 'MergeIntervals',
  setup() {
    const intervalsInput = ref('');
    const mergedIntervals = ref<{ start: number; end: number }[]>([]);
    const errors = ref<string[]>([]);

    // puts form into the initial loading state
    const resetForm = () => {
      errors.value = [];
      mergedIntervals.value = [];
      intervalsInput.value = '';
    };

    // submits form to backend
    const submitForm = async () => {
      errors.value = [];
      mergedIntervals.value = [];

      const intervals = validateIntervals(intervalsInput.value);
      if (errors.value.length > 0) {
        return;
      }

      try {
        const response = await axios.post('http://localhost:8080/merge', { intervals });
        mergedIntervals.value = response.data.result;
      } catch (error) {
        errors.value.push('A server error occurred');
      }
    };

    // validates user input
    const validateIntervals = (input: string): { start: number; end: number }[] | null => {
      // search for invalid input outside of enclosing square brackets
      let nonIntervalMatches = [
        ...input.matchAll(/(\]|^)\s?([^ [\]]+)\s?(\[|$)/g),
        ...input.matchAll(/(\[\s?([^[\]]+)$)/g)
      ];
      for (const nonIntervalMatch of nonIntervalMatches) {
        errors.value.push(`Not an interval: ${nonIntervalMatch[2]}`);
      }

      // validate input inside of enclosing square brackets
      const intervals = input
        .match(/\[([^[\]]+)\]/g)
        ?.map((interval) => {
          const numbers = interval.slice(1, -1).split(',').map(Number);
          if (
            numbers.length !== 2 ||
            isNaN(numbers[0]) ||
            isNaN(numbers[1]) ||
            !Number.isInteger(numbers[0]) ||
            !Number.isInteger(numbers[1]) ||
            numbers[0] > numbers[1]
          ) {
            errors.value.push(`Invalid interval: ${interval}`);
            return null;
          }
          return { start: numbers[0], end: numbers[1] };
        })
        .filter((interval): interval is { start: number; end: number } => interval !== null);

      // no valid intervals found
      if (!intervals) {
        errors.value.push(`No valid interval(s) provided`);
      }

      return intervals && intervals.length > 0 ? intervals : null;
    };

    return {
      intervalsInput,
      mergedIntervals,
      errors,
      submitForm,
      resetForm
    };
  }
});
</script>
