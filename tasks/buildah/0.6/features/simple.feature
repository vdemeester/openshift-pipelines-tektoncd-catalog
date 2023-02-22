Feature: Dockerfile builds with buildah
  Testing different scenario where `buildah` builds Dockerfiles

  Background:
    Given I have a local registry running
    And I have checked out https://github.com/kelseyhightower/nocode in a workspace

  Scenario: push
    When I use the build Task with the following parameters:
      | name  | value                |
      | IMAGE | registry:5000/nocode |
    Then the TaskRun should suceed
    And an image should be published at https://registry:5000/nocode