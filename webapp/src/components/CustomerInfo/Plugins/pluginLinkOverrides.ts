// not all plugins have a homepage url, so this should server as a way to get the plugin ID and still assign it a url.

import {CustomerPluginValues} from '@/types/customers';

const addMissingHomePageURL = (plugin: CustomerPluginValues): CustomerPluginValues => {
    if (plugin.homePageURL) {
        return plugin;
    }

    let homepageURL = '';

    switch (plugin.pluginID) {
    case 'com.mattermost.nps':
        homepageURL = 'https://github.com/mattermost/mattermost-plugin-nps/';
        break;
    case 'com.mattermost.wrangler' :
        homepageURL = 'https://github.com/gabrieljackson/mattermost-plugin-wrangler';
        break;
    case 'com.github.scottleedavis.mattermost-plugin-remind':
        homepageURL = 'https://github.com/scottleedavis/mattermost-plugin-remind';
        break;
    case 'rssfeed':
        homepageURL = 'https://github.com/wbernest/mattermost-plugin-rssfeed';
        break;
    default:
        return plugin;
    }

    return {
        ...plugin,
        homePageURL: homepageURL,
    };
};

export {
    addMissingHomePageURL,
};
