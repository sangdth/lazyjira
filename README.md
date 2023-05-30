# lazyjira

Jira in terminal for lazy people

<img src="https://github.com/sangdth/lazyjira/assets/1083478/05b02f93-273f-4629-8701-12f0d7272d45" width="100%">

🔴 Under development, be careful when using it 🔴

# Getting Started

First, create a configuration file at `~/.config/lazyjira/config.yaml`

```yaml
# Inside ~/.config/lazyjira/config.yaml
username: yourname@email.com
server: https://yourproject.atlassian.net
```

For API token, after generate from [Atlassian](https://id.atlassian.com/manage-profile/security/api-tokens), please add a new record into `Keychain.app`:

- Name: `lazyjira`
- Account: `yourname@email.com`
- Password: Your API token
