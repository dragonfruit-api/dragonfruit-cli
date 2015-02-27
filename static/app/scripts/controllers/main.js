'use strict';

/**
 * @ngdoc function
 * @name blackbookApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the blackbookApp
 */
angular.module('blackbookApp')
  .controller('MainCtrl', function ($scope,vendor) {
    $scope.vendors = vendor.get();

    $scope.showVendor =  function(index){
      console.log($scope.vendors[index]);
    };

  });
