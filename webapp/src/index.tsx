import React from 'react';
import {Store, Action} from 'redux';

import {GlobalState} from '@mattermost/types/store';

import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {manifest} from '@/manifest';

import {PluginRegistry} from '@/types/mattermost-webapp';

import {CustomerRHS} from './app';
import {RHSTitlePlaceholder} from './components/rhsTitle';

const action = `${manifest.id}_received_rhs_action`;
export function setShowRHSAction(showRHSPluginAction: () => void): unknown {
    return {
        type: action,
        showRHSPluginAction,
    };
}

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<string>>) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/

        // RHS Registration
        const {toggleRHSPlugin} = registry.registerRightHandSidebarComponent(CustomerRHS, <RHSTitlePlaceholder/>);
        const boundToggleRHSAction = () => store.dispatch(toggleRHSPlugin);

        // App Bar icon
        if (registry.registerAppBarComponent) {
            // @ts-ignore
            const config = getConfig(store.getState());
            const siteUrl = (config && config.SiteURL) || '';
            const iconURL = `${siteUrl}/plugins/${manifest.id}/public/app-bar-icon.png`;
            registry.registerAppBarComponent({
                iconUrl: iconURL,
                tooltipText: 'Customers',
                action: boundToggleRHSAction,
                supportedProductIds: '*',
            });
        }
    }
}

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
