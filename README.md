# Cleaner

cleaner configuration indicates time to clean data from database. It is defined within key `cleaner`, it is the map where key is the function name and the value is the config for this function.

```json
{
  "cleaner": {
    "ClearSuccessTasks": {
      "cron": "0 0 1 * *",
      "removeAfter": 7776000,
      "skipIfLessThan": 1000,
      "description": "triggered at 00:00 on day-of-month 1"
    }
  }
}
```

Above sample will trigger `ExecClearSuccessTasks` function and start cleaning data every first day of month.
- cron: cron job
- removeAfter: delete data that is created before this time, 7776000 (seconds) is 90 days
- skipIfLessThan: skip cleaning data if there is less than 1000 records left over
- description: describe the task.
