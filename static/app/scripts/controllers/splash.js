'use strict';

/**
 * @ngdoc function
 * @name blackbookApp.controller:SplashCtrl
 * @description
 * # SplashCtrl
 * Controller of the blackbookApp
 */
angular.module('blackbookApp')
  .controller('SplashCtrl', function (authService,$scope,$q) {

    $scope.errorText = $q.when(authService.getLastError);





  });
