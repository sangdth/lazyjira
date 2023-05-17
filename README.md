# lazyjira

Jira in terminal for lazy people

# Getting Started

I will make the automatic generate configuration later, for now please create it:

```yaml
# Inside ~/.config/lazyjira/config.yaml
server: https://yourproject.atlassian.net
login: yourname@email.com
```

For API token, after generate from [Atlassian](https://id.atlassian.com/manage-profile/security/api-tokens), please add a new record into `Keychain.app`:

- Name: `lazyjira`
- Account: `yourname@email.com`
- Password: Your API token

TODO
