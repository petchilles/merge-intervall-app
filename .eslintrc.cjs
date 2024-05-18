/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution');

module.exports = {
  root: true,
  parser: 'vue-eslint-parser',
  extends: [
    'plugin:vue/vue3-essential',
    'eslint:recommended',
    '@vue/eslint-config-typescript',
    '@vue/eslint-config-prettier/skip-formatting',
    '@vue/eslint-config-prettier',
    '@vue/prettier'
  ],
  overrides: [
    {
      files: ['e2e/**/*.{test,spec}.{js,ts,jsx,tsx}'],
      extends: ['plugin:playwright/recommended']
    }
  ],
  parserOptions: {
    ecmaVersion: 'latest'
  },
  rules: {
    'no-console': 1, // Means warning
    'prettier/prettier': 1, // Means warning
    strict: 0, // Means off
    'no-unused-vars': 1,
    'vue/no-unused-components': 1,
    'vue/multi-word-component-names': 0,
    'vue/require-v-for-key': 0
  }
};
