// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.
export const environment = {
  production: false,
  redirectUri: 'http://localhost:4200',
  clientId: '2d929579-7af0-4041-bae1-a18f49858903',
  authority: 'https://login.microsoftonline.com/25d2a676-d0f3-438e-be7b-f55baaf9c7a1',
  serverBaseUrl: 'http://localhost:8080',

};



/*
 * In development mode, to ignore zone related error stack frames such as
 * `zone.run`, `zoneDelegate.invokeTask` for easier debugging, you can
 * import the following file, but please comment it out in production mode
 * because it will have performance impact when throw error
 */
// import 'zone.js/plugins/zone-error';  // Included with Angular CLI.
