import { StackContext, Function, Queue, Api } from "sst/constructs";

export function YtAPIStack({ stack }: StackContext) {
  const redisUrl = process.env.REDIS_URL!;
  const ytApiKey = process.env.API_KEY!;
  const dbUrl = process.env.DATABASE_URL!;

  const api = new Api(stack, "api", {
    customDomain: "task.volt.place",
    routes: {
      "GET /videos": "cmd/list_videos/main.go",
    },
    defaults: {
      function: {
        environment: {
          REDIS_URL: redisUrl,
          API_KEY: ytApiKey,
          DATABASE_URL: dbUrl,
        },
      },
    },
  });

  const saveVideos = new Function(stack, "saveVideos", {
    handler: "cmd/save_videos/main.go",
    environment: {
      REDIS_URL: redisUrl,
      DATABASE_URL: dbUrl,
      API_KEY: ytApiKey,
    },
  });

  stack.addOutputs({
    ApiEndpoint: api.url,
  });
}
