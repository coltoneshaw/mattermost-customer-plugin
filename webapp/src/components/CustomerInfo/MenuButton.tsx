import {ActionIcon, Menu, createStyles} from '@mantine/core';
import React from 'react';

import {MenuIcon} from '@mattermost/compass-icons/components';
import {generatePath, useHistory} from 'react-router-dom';

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

const menuItems = [
    {
        value: '',
        label: 'Information',
    },
    {
        value: '/packet',
        label: 'Support Packet',
    },
    {
        value: '/config',
        label: 'Config',
    },
    {
        value: '/plugins',
        label: 'Plugins',
    },
];

type Params = {
    id: string
}
const MenuButton = ({
    id,
}: Params) => {
    const {classes} = useStyles();
    const history = useHistory();

    return (
        <Menu
            classNames={classes}
            position='bottom-end'
            offset={10}
        >
            <Menu.Target>
                <ActionIcon
                    sx={{
                        height: '32px',
                        width: '32px',
                        color: 'rgba(var(--center-channel-color-rgb), 0.72)',
                        '&:hover': {
                            backgroundColor: 'rgba(var(--center-channel-color-rgb), 0.08)',
                            color: 'rgba(var(--center-channel-color-rgb), 0.72)',
                        },
                        marginLeft: 'auto',
                    }}
                >
                    <MenuIcon
                        size={18}
                    />
                </ActionIcon>
            </Menu.Target>
            <Menu.Dropdown>
                <Menu.Item
                    style={{
                        color: 'rgba(var(--center-channel-color-rgb),0.56)',
                        fontWeight: 600,
                        textTransform: 'uppercase',
                        fontSize: '12px',
                    }}
                    disabled={true}
                >
                    {'Menu'}
                </Menu.Item>
                {
                    menuItems.map((item) => {
                        const path = generatePath('/customers/:id' + item.value, {
                            id,
                        });
                        return (
                            <Menu.Item
                                key={item.value}
                                onClick={() => {
                                    history.push(path);
                                }}
                                style={{
                                    fontSize: '14px',
                                    color: 'var(--center-channel-color-rgb)',
                                }}
                            >
                                {item.label}
                            </Menu.Item>
                        );
                    })
                }
            </Menu.Dropdown>
        </Menu>
    );
};

export {
    MenuButton,
};
