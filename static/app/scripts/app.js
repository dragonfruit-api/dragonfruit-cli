'use strict';

/**
 * @ngdoc overview
 * @name blackbookApp
 * @description
 * # blackbookApp
 *
 * Main module of the application.
 */
angular
  .module('blackbookApp', [
    'ngAnimate',
    'ngResource',
    'ngRoute',
    'ngSanitize',
    'LocalStorageModule'
  ])
  .config(function ($routeProvider, localStorageServiceProvider) {
    localStorageServiceProvider.setPrefix('76f6e66d-b739-41a2-9af5-a982bcad4998');
    var authorizer = ['$q', 'authService', function($q, authService) {

      var userInfo = authService.getUserInfo();
      if (userInfo) {
        return $q.when(userInfo);
      } else {
        return $q.reject({ authenticated: false });
      }
    }];
    $routeProvider
      .when('/loginredirect', {
        controller: 'LoginRedirectCtrl',
        templateUrl: 'views/splash.html'
      })    
      .when('/login', {
        templateUrl: 'views/splash.html',
        controller: 'SplashCtrl'
      })
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl',
        resolve: {
           auth: authorizer
        },
      })
      .when('/logout', {
        templateUrl: 'views/logout.html',
        controller: 'LogoutCtrl'
      })
      .otherwise({
        redirectTo: '/loginredirect'
      });
  }).run(
    ['$rootScope','$location','$anchorScroll', 
    function($rootScope,$location,$anchorScroll){
    
      $rootScope.$on('$routeChangeSuccess', function(event, toState) {
        console.log(event, toState);
        $rootScope.authstate = true;
        //$location.hash('top');
        //$anchorScroll();

      });
   
      $rootScope.$on('$routeChangeError', function(event, current, previous, eventObj) {
        console.log('login error: ',event,current, previous,eventObj);
        if (eventObj.authenticated === false) {
          $rootScope.authstate = false;
          $location.path('/login');
        }
      });
    }]
  );

