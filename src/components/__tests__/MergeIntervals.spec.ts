import { mount, flushPromises } from '@vue/test-utils';
import MergeIntervals from '@/components/MergeIntervals.vue';
import { describe, expect, it, vi } from 'vitest';
import { createVuetify } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';
import axios from 'axios';

const vuetify = createVuetify({ components, directives });

// always mock response data to test component independently
const mockResponse = {
  data: {
    result: [
      [2, 23],
      [25, 30]
    ]
  }
};

// tests component MergeIntervals.vue with 100% coverage
describe('MergeIntervals.vue', () => {
  const wrapper = mount(MergeIntervals, {
    global: {
      plugins: [vuetify]
    }
  });

  it('renders as expected', () => {
    expect(wrapper.html()).toMatchFileSnapshot('./__snapshots__/mergeIntervals.html');
  });

  it('handles valid input data correctly', async () => {
    // mock response
    vi.spyOn(axios, 'post').mockResolvedValue(mockResponse);

    // Data can be entered into the pprovided textarea?
    const mockInput = '[25,30] [2,19] [14,23] [4,8]';
    const textarea = wrapper.find('textarea');
    await textarea.setValue(mockInput);
    expect(textarea.element.value).toBe(mockInput);

    // Form can be submitted?
    await wrapper.find('form').trigger('submit.prevent');
    await flushPromises();
    expect(axios.post).toHaveBeenCalled();

    // Component is updated with mocked result?
    const mergedIntervals = wrapper.vm.mergedIntervals;
    expect(mergedIntervals).toEqual(mockResponse.data.result);
  });

  it('resets the form correctly', async () => {
    // check if reset button exists
    const resetButton = wrapper.find('#reset');
    expect(resetButton.exists()).toBe(true);

    // check if pressing reset button resets the form
    await resetButton.trigger('click');
    await flushPromises();

    // expect result header to be invisible
    const resultHeader = wrapper.find('#result_header');
    expect(resultHeader.exists()).toBe(false);

    // expect no results from the last submission to be displayed
    const results = wrapper.find('#results');
    expect(results.exists()).toBe(false);
  });

  it('handles invalid iput data correctly', async () => {
    // enter invalid intervals string
    const mockInput = 'abc[25,30] [2,19] [14, 23a] def [4,8 ghi';
    const textarea = wrapper.find('textarea');
    await textarea.setValue(mockInput);
    expect(textarea.element.value).toBe(mockInput);

    // check if submit button [merge] exists
    const mergeButton = wrapper.find('#merge');
    expect(mergeButton.exists()).toBe(true);

    // press submit button
    await mergeButton.trigger('submit.prevent');
    await flushPromises();

    // check if info message appears
    const infoElement = wrapper.find('#info');
    expect(infoElement.exists()).toBe(true);

    // check validation message
    const errorElement = wrapper.find('#error');
    expect(errorElement.exists()).toBe(true);
    expect(errorElement.html()).toMatchFileSnapshot('./__snapshots__/validation.message_1.html');
  });

  it('handles empty form submission correctly', async () => {
    // reset form
    wrapper.find('#reset').trigger('click');

    // submit empty form
    await wrapper.find('form').trigger('submit.prevent');
    await flushPromises();

    // check (again) if info message appears
    const infoElement = wrapper.find('#info');
    expect(infoElement.exists()).toBe(true);

    // check for correct validation message
    const errorElement = wrapper.find('#error');
    expect(errorElement.exists()).toBe(true);
    expect(errorElement.html()).toMatchFileSnapshot('./__snapshots__/validation.message_2.html');
  });

  it('handles response failure', async () => {
    vi.spyOn(axios, 'post').mockResolvedValue('net::ERR_CONNECTION_REFUSED');

    // mock form with valid data and submit
    const mockInput = '[5,30] [2,4] [4,28]';
    const textarea = wrapper.find('textarea');
    await textarea.setValue(mockInput);
    await wrapper.find('form').trigger('submit.prevent');
    await flushPromises();
    expect(axios.post).toHaveBeenCalled();

    // check validation message
    const errorElement = wrapper.find('#error');
    expect(errorElement.exists()).toBe(true);
    expect(errorElement.html()).toMatchFileSnapshot('./__snapshots__/server.error.message.html');
  });
});
