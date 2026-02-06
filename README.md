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

### 2. Create Slack Bot

1. Go to https://api.slack.com/apps/
2. Click "Create New App" → "From manifest"
3. Paste this manifest in JSON editor:
```json
{
    "display_information": {
        "name": "GitHub PR Notifier",
        "description": "Personal GitHub PR notification bot - sends low-noise Slack DMs for PR activity"
    },
    "features": {
        "bot_user": {
            "display_name": "GitHub PR Bot",
            "always_online": true
        }
    },
    "oauth_config": {
        "scopes": {
            "bot": [
                "chat:write",
                "im:write"
            ]
        }
    },
    "settings": {
        "org_deploy_enabled": false,
        "socket_mode_enabled": false,
        "is_hosted": false,
        "token_rotation_enabled": false
    }
}
```
4. **Workspace**: Select your workspace
5. Copy the **Bot User OAuth Token** (starts with `xoxb-`)
6. Copy your **Slack User ID** (format: `U1234567890`)
