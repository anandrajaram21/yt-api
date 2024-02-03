import { SSTConfig } from "sst";
import { Api } from "sst/constructs";
import { YtAPIStack } from "./stacks/YtAPIStack";

export default {
  config(_input) {
    return {
      name: "yt-api",
      region: "ap-south-1",
    };
  },
  stacks(app) {
    app.setDefaultFunctionProps({
      runtime: "go",
    });
    app.stack(YtAPIStack);
  },
} satisfies SSTConfig;
