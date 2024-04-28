import { SSTConfig } from 'sst'
import { API } from './stacks/MyStack'
import { Astro } from './stacks/Astro';

export default {
  config(_input) {
    return {
      name: "railopedia",
      region: "eu-west-2",
    };
  },
  stacks(app) {
    app.stack(API).stack(Astro)  
  }
} satisfies SSTConfig;
