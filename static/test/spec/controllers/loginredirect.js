'use strict';

describe('Controller: LoginredirectctrlCtrl', function () {

  // load the controller's module
  beforeEach(module('blackbookApp'));

  var LoginredirectctrlCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    LoginredirectctrlCtrl = $controller('LoginredirectctrlCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
