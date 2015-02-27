'use strict';

/**
 * @ngdoc function
 * @name blackbookApp.controller:LoginredirectctrlCtrl
 * @description
 * # LoginredirectctrlCtrl
 * Controller of the blackbookApp
 */
angular.module('blackbookApp')
  .controller('LoginRedirectCtrl', function (authService, $location) {
    var userInfo = authService.getUserInfo();
    console.log('user info:',userInfo);
    if(userInfo !== null) {
      $location.path('/');
    } 

    var handleError = function(error){
      authService.setLastError(error);
      $location.path('/login');
    };

    var handleSuccess = function(result){
      console.log('result: ',result);
      $location.path('/');
    };    

    authService.login().then(handleSuccess,handleError);
  });
