// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.
export const environment = {
  production: false,
  redirectUri: "http://localhost:4200",
  clientId: "50afc41f-4647-4787-9d61-6c8bec34091c",
  authority:
    "https://login.microsoftonline.com/5806938e-ea7d-4345-85fb-6239156b78d6",
  serverBaseUrl: "http://localhost:8080",
};

/*
 * In development mode, to ignore zone related error stack frames such as
 * `zone.run`, `zoneDelegate.invokeTask` for easier debugging, you can
 * import the following file, but please comment it out in production mode
 * because it will have performance impact when throw error
 */
// import 'zone.js/plugins/zone-error';  // Included with Angular CLI.
