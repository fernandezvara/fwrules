(function() {
  var app = angular.module('fwrules', [
    'ui.router',
    'ngResource',
    'ui.bootstrap',
    'xeditable'
  ]);

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

  app.factory('Servers', function($resource) {
    return $resource('/api/test/machines/:id');
  });

  app.factory('RuleSets', function($resource) {
    return $resource('/api/test/rulesets/:name', {}, {
      query: {
        method: 'GET',
        params: {
          name: ''
        },
        isArray: true
      },
      get: {
        method: 'GET',
        params: {
          name: '@name'
        }
      },
      post: {
        method: 'POST'
      },
      update: {
        method: 'PUT',
        params: {
          name: '@name'
        }
      },
      remove: {
        method: 'DELETE',
        params: {
          name: '@name'
        }
      }
    });
  });

  app.run(function(editableOptions) {
    editableOptions.theme = 'bs3';
  });

  app.config(function($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise("/");

    $stateProvider
      .state('home', {
        url: '/',
        templateUrl: '/app/views/pageHome.html'
      })
      .state('servers', {
        url: '/servers',
        templateUrl: 'app/views/pageServers.html',
        controller: ['$scope', 'Servers',
          function($scope, Servers) {
            Servers.query(function(data) {
              $scope.allServers = data;
            });
          }
        ]
      })
      .state('rulesets', {
        url: '/rulesets',
        templateUrl: 'app/views/pageRulesets.html',
        controller: ['$scope', 'RuleSets',
          function($scope, RuleSets) {
            RuleSets.query(function(data) {
              $scope.allRulesets = data;
            });
          }
        ]
      })
      .state('rulesets.detail', {
        url: '/{name}',
        views: {
          'detail': {
            controller: ['$scope', '$state', '$stateParams', 'RuleSets',
              '$modal',
              function($scope, $state, $stateParams, RuleSets, $modal) {
                RuleSets.get({
                  name: $stateParams.name
                }, function(data) {
                  $scope.ruleset = data;
                });

                $scope.updateRuleset = function() {
                  RuleSets.update($scope.ruleset);
                  console.log('updated!');
                };

                $scope.removeModal = function() {
                  RuleSets.delete({
                    name: $scope.ruleset.name
                  });
                  $state.go('^');
                  console.log('user said true');
                };
              }
            ],
            templateUrl: 'app/views/pageRulesetsDetail.html'
          }
        }
      });
  });

  app.directive('ngConfirmClick', [
    function() {
      return {
        link: function(scope, element, attr) {
          var msg = attr.ngConfirmClick || "Are you sure?";
          var clickAction = attr.confirmedClick;
          element.bind('click', function(event) {
            if (window.confirm(msg)) {
              scope.$eval(clickAction);
            }
          });
        }
      };
    }
  ]);

})();
console.log('App loaded...');
