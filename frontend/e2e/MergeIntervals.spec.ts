import { test, expect } from '@playwright/test';

test('App returns merged intervals on correct input request', async ({ page }) => {
  // Navigate to the start page
  await page.goto('/');

  // Enter test data (intervals)
  const testData = '[25,30] [2,19] [14,23] [4,8]';
  await page.fill('textarea[id="input"]', testData);

  // Click the merge button
  await page.click('button[id="merge"]');

  // Wait until the result text area is visible
  const resultTextarea = page.locator('textarea[id="result"]');
  await resultTextarea.waitFor({ state: 'visible' });

  // Extract calculated result coming from the golang backend
  const mergedIntervalsText = await resultTextarea.evaluate(
    (node) => (node as HTMLTextAreaElement).value
  );

  // Expected result
  const expectedMergedIntervals = '[2,23][25,30]';

  // Check whether the text of the result textarea corresponds to the expected result
  expect(mergedIntervalsText?.trim()).toBe(expectedMergedIntervals);
});
