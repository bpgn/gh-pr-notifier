### 1. Create GitHub App

1. Go to https://github.com/settings/apps/new
2. Configure:
   - **App name**: `GitHub PR Notifier` (or your choice)
   - **Homepage URL**: `https://yourserver.com` (where you'll run this service)
   - **Webhook URL**: `https://yourserver.com/webhook`
   - **Webhook secret**: Generate a strong random secret (save this)
3. **Permissions**:
   - Repository → `Pull requests`: Read-only
4. **Subscribe to events**:
   - Pull request
   - Pull request review
   - Pull request review comment
5. Click "Create GitHub App"
6. Generate private key (for API access, optional)
7. Install app on your repositories
