'use strict';

/**
 * @ngdoc service
 * @name staticApp.authService
 * @description
 * # authService
 * Factory for handling authorization.
 */
angular.module('blackbookApp')
  .factory('authService', function($http, $q, localStorageService, $location) {
  var userInfo,
    lastError;

  var getUserInfo = function() {
    return localStorageService.get('userInfo');
  };
 
  var login = function() {
    var deferred = $q.defer();
 
    $http.post('/checkauth').then(function(result) {
      console.log(result);
      // if the user is not logged in reject the promise.
      if (result.data.isLoggedIn === false) {
        deferred.reject('User is not logged in.');
        return;
      }

      // otherwise save the user to local storage so we don't have to 
      userInfo = result.data.user;
      localStorageService.set('userInfo',userInfo);
      lastError = '';
      deferred.resolve(userInfo);

    }, function(error) {
      console.log('test login error:',error);
      // I like turtles. 
      deferred.reject(error);
    });
 
    return deferred.promise;
  };

  var getLastError = function(){
    return lastError;
  };

  var setLastError  = function(error){
    lastError = error;
  };

  function logout(){
    localStorageService.remove('userInfo');
    $location.path('/logout');
  }
 
  return {
    login: login,
    getUserInfo: getUserInfo,
    logout: logout,
    getLastError: getLastError,
    setError: setLastError
  };
});
