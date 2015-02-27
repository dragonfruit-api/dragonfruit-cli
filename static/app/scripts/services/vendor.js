'use strict';

/**
 * @ngdoc service
 * @name blackbookApp.vendor
 * @description
 * # vendor
 * Service in the blackbookApp.
 */
angular.module('blackbookApp')
  .factory('vendor', function () {
    // AngularJS will instantiate a singleton by calling "new" on this function
    var vendorlist = [
      {
        name: '1 Over 0',
        url: 'http://1over0.com/',
        locations: ['New York City'],
        contacts: [
          {
            name: 'Paul Williams',
            title: 'President',
            phone: '646-564-3599',
            email: 'pwilliams@1over0.com'
          }
        ],
        size: 'small',
        ideoReferences: [
          {
            name: 'Peter Olson',
            email: 'polson@ideo.com'
          }
        ],
        image: null,
        competencies: [
          'agile',
          'web'
        ],
        technologies: [
          '.NET',
          'PHP',
          'MSSQL'
        ],
        ideoProjects:[
          {
            client: 'Walgreens',
            codename: 'Orange',
            sector: 'Government',
            website: 'http://ideorama.ideo.com/project-view.php?project=walgreen-co-walgreens----future-pharmacy&view=index&projectId=2305#page-0'
          }
        ],
        engagementModel: '200/hr/employee, minimum 2 weeks',
        notes: [
          {
            author: 'Peter',
            comment: ''
          }
        ],
        status: 'active',
        ideoConfidence: 'preferred'
      }
    ];

    return {
      get: function(){return vendorlist;}
    }; 
  });
