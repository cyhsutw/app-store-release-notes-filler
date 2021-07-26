# Release Notes Filler for App Store Connect

A simple web app that fills release notes on App Store Connect by using a given localization key from Lokalise.

![Demo](/docs/demo.gif)

## Configuration

### Describe apps
Create an `apps.yml` from template `apps.yml.exmaple`:

```yaml
- name: 'My Cool App'
  app_id: 'some app store id'
  lokalise_project_id: 'some lokalise project id'
  icon_url: 'some image url'

- name: 'Another Cool App'
  app_id: 'another app store id'
  lokalise_project_id: 'another lokalise project id'
  icon_url: 'another image url'
```

- `name`: The name of the app, visible on the web UI for this web app.
- `app_id`: You can find this on the App Store Connect page of your app. Under General > App Information > Apple ID.
- `lokalise_project_id`: You can find this on the Lokalise project page.
- `icon_url`: The url of the app icon, you may place the icon under `/images` and set this value as something like `/images/xxx.png`.

### Setup API keys

Create an `.env` from template `.env.example`:

```
LOKALISE_API_KEY=xxx
APP_STORE_CONNECT_API_KEY_ID=yyy
APP_STORE_CONNECT_API_ISSUER_ID=xxx
APP_STORE_CONNECT_API_PRIVATE_KEY="single line private key"
```

- `LOKALISE_API_KEY`: a Lokalise API key that can access the projects configured in `apps.yml`
- `APP_STORE_CONNECT_API_KEY_ID`: the key ID of the App Store Connect API key
- `APP_STORE_CONNECT_API_ISSUER_ID`: the issuer ID of the App Store Connect API key
- `APP_STORE_CONNECT_API_PRIVATE_KEY`: the private key of App Store Connect API key.

  Format it as a single-line string, wrap everything in a pair of double quotes, and, replace newlines with `\n`.

  For example:

  ```
  APP_STORE_CONNECT_API_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\n6WHB/XNDdnq4U291fnFK+ytueycdogynUXK9Xvt1x/F+4VN9eAsh4KARKd/kI+aJ\nZ0FjGwOP9QOUmpu8u1H5pWfXg0NPUkzUPj4X8sIA11iwV5kClOVV3WfANJpKtGx2\nWLTrX9SWQtFdP+PDJQjcINRt7Chqh028Q3rCeodLmR6KDmcCEsB9pgfj36JYmb0F\neGTTbdJG\n-----END PRIVATE KEY-----"
  ```

## Deployment

### Install Docker and Docker Compose

Visit https://docs.docker.com/get-docker/ for instructions.

### Boot the app

```bash
$ docker-compose up
```

And visit `http://localhost:48745`


# Development

- Go `1.16`
- Docker
- Docker Compose
