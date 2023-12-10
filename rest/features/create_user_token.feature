Feature: User token
    In order to call endpoints behind authorization it is required to use a user token
    As unuathorized user
    Needs to create a token by providing credentials
    And then it is possible to call endpoints with authorization

    Scenario: Create a new token by providing correct credentials
        Given valid username and password
        When call /user/token endpoint
        Then user token is returned
        
    Scenario: Get user information with a valid JWT token
        Given valid user token
        When call /user/info endpoint
        Then user information is returned

    Scenario: Get user information with an invalid JWT token
        Given invalid user token
        When call /user/info endpoint
        Then 401 status is returned