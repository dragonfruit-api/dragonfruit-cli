'use strict';

/**
 * @ngdoc function
 * @name blackbookApp.controller:LogoutCtrl
 * @description
 * # LogoutCtrl
 * Controller of the blackbookApp
 */
angular.module('blackbookApp')
  .controller('LogoutCtrl', function (authService) {
    authService.logout();
  });
