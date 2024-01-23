import React from 'react';

import {CustomerPluginValues} from '@/types/customers';
import {CenteredText} from '@/components/CenteredText';
import {Container} from '@/components/Group';
import {dynamicSort} from '@/helpers/sort';

import {MenuButton} from './PluginMenu';
import {addMissingHomePageURL} from './pluginLinkOverrides';

type PluginCardParams = {
    plugin: CustomerPluginValues;
    changeState: (pluginID: string, isActive: boolean) => void;
    deletePlugin: (pluginID: string) => void;
}
const PluginCard = ({
    plugin,
    changeState,
    deletePlugin,
}: PluginCardParams) => {
    return (
        <div
            style={{
                border: '1px solid rgba(var(--center-channel-color-rgb), 0.16)',
                background: 'var(--center-channel-bg)',
                color: 'var(--center-channel-color)',

                // cursor: 'pointer',
                height: '6.4em',
                width: '100%',
                boxShadow: 'rgba(0,0,0,0.08) 0px 2px 3px 0px',
                padding: '1em',
                display: 'flex',
                flexDirection: 'row',
                justifyContent: 'space-between',
            }}
        >
            <div
                style={{
                    display: 'flex',
                    flexDirection: 'column',

                    maxWidth: 'calc(100% - 60px)',
                    paddingRight: '1em',
                    justifyContent: 'space-between',
                }}
            >
                <span
                    style={{
                        whiteSpace: 'nowrap',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                    }}
                >
                    <strong>{plugin.name}</strong>
                    {` (${plugin.version})`}
                </span>
                <span
                    style={{
                        whiteSpace: 'nowrap',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                    }}
                >
                    {plugin.pluginID}
                </span>
            </div>
            <div
                style={{
                    display: 'flex',
                    flexDirection: 'column',
                    justifyContent: 'space-between',
                    alignContent: 'flex-end',
                }}
            >
                <MenuButton
                    plugin={plugin}
                    changeState={changeState}
                    deletePlugin={deletePlugin}
                />
                <span
                    style={{
                        padding: '4px 8px',
                        borderRadius: '4px',
                        fontSize: '12px',
                        background: plugin.isActive ? 'var(--online-indicator)' : 'var(--away-indicator)',
                        color: 'var(--sys-button-color)',
                        fontWeight: 600,
                    }}
                >
                    {plugin.isActive ? 'Active' : 'Inactive'}
                </span>
            </div>

        </div>
    );
};
type Params = {
    plugins: CustomerPluginValues[] | null;
    changeState: (pluginID: string, isActive: boolean) => void;
    deletePlugin: (pluginID: string) => void;
}
const CustomerInfoPlugins = ({
    plugins,
    changeState,
    deletePlugin,
}: Params) => {
    if (!plugins || plugins.length === 0) {
        return (
            <CenteredText
                message={'No plugin information found.'}
            />
        );
    }
    return (
        <Container>
            {
                plugins.sort(dynamicSort('name')).map((p) => {
                    return (
                        <PluginCard
                            key={p.pluginID}
                            plugin={addMissingHomePageURL(p)}
                            changeState={changeState}
                            deletePlugin={deletePlugin}
                        />
                    );
                })
            }
        </Container>
    );
};

export {
    CustomerInfoPlugins,
};
