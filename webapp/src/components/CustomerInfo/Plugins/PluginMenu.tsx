import {ActionIcon, Menu, createStyles} from '@mantine/core';
import React from 'react';

import {DotsHorizontalIcon} from '@mattermost/compass-icons/components';

import {CustomerPluginValues} from '@/types/customers';

const useStyles = createStyles(() => ({
    dropdown: {
        minWidth: '225px',
        maxWidth: '270px',
        boxShadow: '0 6px 12px rgba(0, 0, 0, 0.175)',
        borderRadius: '4px',
        border: '1px solid rgba(0, 0, 0, 0.15)',
        backgroundColor: 'var(--center-channel-bg)',
        borderColor: 'rgba(221, 223, 228, 0.2)',
    },
    item: {
        height: '34px',
        padding: '1px 16px',
        '&:hover': {
            backgroundColor: 'rgba(221, 223, 228, 0.1)',
        },
    },
    itemLabel: {
        lineHeight: '22px',
        padding: '5px 0',
        whiteSpace: 'normal',
        display: 'flex',
        justifyContent: 'space-between',
    },
}));

type Params = {
    plugin: CustomerPluginValues;
    changeState: (pluginID: string, isActive: boolean) => void;
    deletePlugin: (pluginID: string) => void;
}
const MenuButton = ({
    plugin,
    changeState,
    deletePlugin,
}: Params) => {
    const {classes} = useStyles();

    return (
        <Menu
            classNames={classes}
            position='bottom-end'
            offset={5}
        >
            <Menu.Target>
                <ActionIcon
                    sx={{
                        height: '28px',
                        width: '28px',
                        color: 'rgba(var(--center-channel-color-rgb), 0.72)',
                        '&:hover': {
                            backgroundColor: 'rgba(var(--center-channel-color-rgb), 0.08)',
                            color: 'rgba(var(--center-channel-color-rgb), 0.72)',
                        },
                        marginLeft: 'auto',
                    }}
                >
                    <DotsHorizontalIcon
                        size={16}
                    />
                </ActionIcon>
            </Menu.Target>
            <Menu.Dropdown>
                <Menu.Item
                    style={{
                        fontSize: '14px',
                        color: 'var(--center-channel-color-rgb)',
                    }}
                    onClick={() => {
                        changeState(plugin.pluginID, !plugin.isActive);
                    }}
                >
                    {'Mark ' + (plugin.isActive ? 'Inactive' : 'Active')}
                </Menu.Item>
                {plugin.homePageURL && (
                    <Menu.Item
                        style={{
                            fontSize: '14px',
                            color: 'var(--center-channel-color-rgb)',
                        }}

                        onClick={() => {
                            window.open(plugin.homePageURL, '_blank');
                        }}
                    >
                        {'Open Plugin Homepage'}
                    </Menu.Item>
                )}
                <Menu.Item
                    style={{
                        fontSize: '14px',
                        color: 'var(--center-channel-color-rgb)',
                    }}
                    onClick={() => {
                        deletePlugin(plugin.pluginID);
                    }}
                >
                    {'Delete'}
                </Menu.Item>
            </Menu.Dropdown>
        </Menu>
    );
};

export {
    MenuButton,
};
