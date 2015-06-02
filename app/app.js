(function() {
  var app = angular.module('fwrules', ['ui.router', 'restangular']);

  app.directive('mainMenu', function() {
    return {
      restrict: 'E',
      templateUrl: "/app/views/mainMenu.html"
    };
  });

  app.run(
    ['$rootScope', '$state', '$stateParams',
      function($rootScope, $state, $stateParams) {
        $rootScope.$state = $state;
        $rootScope.$stateParams = $stateParams;
      }
    ]
  );

  app.config(function($stateProvider, $urlRouterProvider, RestangularProvider) {
    $urlRouterProvider.otherwise("/");

    $stateProvider
      .state('home', {
        url: '/',
        templateUrl: '/app/views/pageHome.html'
      })
      .state('servers', {
        url: '/servers',
        templateUrl: 'app/views/pageServers.html'
      })
      .state('rulesets', {
        url: '/rulesets',
        templateUrl: 'app/views/pageRulesets.html',
        controller: ['$scope', 'Restangular',
          function($scope, Restangular) {
            var apiRulesets = Restangular.all('rulesets');

            apiRulesets.getList().then(function(rulesets) {
              $scope.allRulesets = rulesets;
            });
          }
        ]
      })
      .state('rulesets.detail', {
        url: '/{name}',
        views: {
          'detail': {
            controller: ['$scope', '$stateParams', 'Restangular',
              function($scope, $stateParams, Restangular) {
                $scope.ruleset = Restangular.one('rulesets',
                  $stateParams.name).get();

                //utils.findById($scope.allRulesets,
                //  $stateParams.name);
              }
            ],
            templateUrl: 'app/views/pageRulesetsDetail.html'
          }
        }
      });

    RestangularProvider.setBaseUrl("/api/test");
    RestangularProvider
      .setRestangularFields({
        id: 'name'
      });
  });

  app.controller('ServersController', function($scope, Restangular) {
    var apiServers = Restangular.all('machines');

    apiServers.getList().then(function(servers) {
      $scope.allServers = servers;
    });
  });

  // app.controller('RuleSetsController', function($scope, Restangular) {
  //   var apiRulesets = Restangular.all('rulesets');
  //
  //   apiRulesets.getList().then(function(rulesets) {
  //     $scope.allRulesets = rulesets;
  //   });
  //
  // });

})();
console.log('File loaded...');
