import {ActionIcon, Menu, createStyles} from '@mantine/core';
import React, {Dispatch} from 'react';

import {CheckIcon, FilterVariantIcon} from '@mattermost/compass-icons/components';

import {CustomerSortOptions} from '@/types/customers';

const useStyles = createStyles(() => ({
    dropdown: {
        minWidth: '225px',
        maxWidth: '270px',
        boxShadow: '0 6px 12px rgba(0, 0, 0, 0.175)',
        borderRadius: '4px',
        border: '1px solid rgba(0, 0, 0, 0.15)',
        backgroundColor: 'var(--center-channel-bg)',
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

type FilterParams = {
    sortBy: CustomerSortOptions;
    setSortBy: Dispatch<React.SetStateAction<CustomerSortOptions>>

}

const menuItems = [
    {
        value: CustomerSortOptions.SortByType,
        label: 'Type',
    },
    {
        value: CustomerSortOptions.SortByName,
        label: 'Name',
    },
    {
        value: CustomerSortOptions.SortByAE,
        label: 'AE',
    },
    {
        value: CustomerSortOptions.SortByCSM,
        label: 'CSM',
    },
    {
        value: CustomerSortOptions.SortByTAM,
        label: 'TAM',
    },
    {
        value: CustomerSortOptions.SortByLicensedTo,
        label: 'Licensed To',
    },
    {
        value: CustomerSortOptions.SortBySiteURL,
        label: 'Site URL',
    },
];

const FilterButton = ({
    sortBy,
    setSortBy,
}: FilterParams) => {
    const {classes} = useStyles();

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
                    }}
                >
                    {sortBy !== CustomerSortOptions.Default && (
                        <i
                            style={{
                                position: 'absolute',
                                top: '6px',
                                right: '6px',
                                width: '10px',
                                height: '10px',
                                border: '2px solid var(--center-channel-bg)',
                                backgroundColor: 'var(--button-bg)',
                                borderRadius: '50%',
                            }}
                        />
                    )
                    }
                    <FilterVariantIcon
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
                    {'Sort Customers By'}
                </Menu.Item>
                {
                    menuItems.map((item) => {
                        return (
                            <Menu.Item
                                key={item.value}
                                onClick={() => {
                                    if (item.value === sortBy) {
                                        setSortBy(CustomerSortOptions.Default);
                                        return;
                                    }
                                    setSortBy(item.value);
                                }}
                                style={{
                                    fontSize: '14px',
                                    color: 'var(--center-channel-color-rgb)',

                                }}
                            >
                                {item.label}
                                {item.value === sortBy && (
                                    <CheckIcon
                                        size={18}
                                    />
                                )}
                            </Menu.Item>
                        );
                    })
                }
            </Menu.Dropdown>

        </Menu>
    );
};

export {
    FilterButton,
};
