'use strict';

/**
 * @ngdoc function
 * @name blackbookApp.controller:AboutCtrl
 * @description
 * # AboutCtrl
 * Controller of the blackbookApp
 */
angular.module('blackbookApp')
  .controller('AboutCtrl', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
  });
