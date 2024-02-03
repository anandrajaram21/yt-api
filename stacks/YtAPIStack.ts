import { StackContext, Function, Queue, Api, Cron } from "sst/constructs";

export function YtAPIStack({ stack }: StackContext) {
  const redisUrl = process.env.REDIS_URL!;
  const ytApiKey = process.env.API_KEY!;
  const dbUrl = process.env.DATABASE_URL!;
  const customDomain = process.env.CUSTOM_DOMAIN;

  const api = new Api(stack, "api", {
    customDomain,
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

  const cron = new Cron(stack, "cron", {
    schedule: "rate(1 minute)",
    job: saveVideos,
  });

  stack.addOutputs({
    ApiEndpoint: api.url,
  });
}
