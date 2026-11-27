const { app } = require('@azure/functions');

app.http('FuncFromCli', {
    methods: ['GET', 'POST'],
    authLevel: 'anonymous',
    handler: async (request, context) => {
        context.log(`Http function processed request for url "${request.url}"`);

        const name = request.query.get('name') || await request.text() || 'world';

/* For easier development and debugging, CORS restrictions have been disabled by setting Access-Control-Allow-Origin: *. 
Once the application is running successfully, ensure CORS is re-enabled and properly configured to secure the application. */
        return {
        body: `Hello ${name}, this is coming from Function App!`,
        headers: {
            "Content-Type": "application/json",
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Methods": "GET, POST, OPTIONS",
            "Access-Control-Allow-Headers": "Content-Type",
          }
    }
    }
});
