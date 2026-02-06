I want to create a small Github PR notifier, that will receive webhook from a Github App and send Slack message to me.

INSTRUCTIONS:
- small iterations, start very simple, test and add stuff
- keep a small and comprehensive README.md
- keep code as simple as possible and don't write too much lines and complex code
- keep safe endpoint and do not introduce regression or security breaches
- keep packages small and easy to understand
- avoid too many dependencies
- avoid writing too much code without testing with me first
- always run tests at the end of an iteration

## Iterations

1. [x] Small golang server that listens on
    - GET /healthz for healthcheck
    - POST /webhook
2. [x] Parse payload from webhook object, events are of type:
    - pull_request_review
    - pull_request_review_comment
3. [x] Notify (dry-run for now using a log) when one of MY open PRs receives:
   - a new review (pull_request_review)
   - a new review comment (pull_request_review_comment)
   - Filter out own actions (no self-notifications)

## Next iterations
4. [ ] Verify webhook signature (GITHUB_WEBHOOK_SECRET)
5. [ ] Send Slack notifications instead of logs

Let's build this together
