// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.
export const environment = {
  production: false,
  redirectUri: 'http://localhost:4200',
  clientId: 'c6d06e75-6f10-42a0-ab0e-55e9bb8c64fe',
  authority: 'https://login.microsoftonline.com/4cab010c-fa07-44f7-bc69-561184a9fb8e',
  serveurBaseUrl: 'http://koki2.com:4200'
};

/*
 * In development mode, to ignore zone related error stack frames such as
 * `zone.run`, `zoneDelegate.invokeTask` for easier debugging, you can
 * import the following file, but please comment it out in production mode
 * because it will have performance impact when throw error
 */
// import 'zone.js/plugins/zone-error';  // Included with Angular CLI.
