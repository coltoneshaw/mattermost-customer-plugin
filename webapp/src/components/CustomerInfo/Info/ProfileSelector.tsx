import {Avatar, Group, Select, Text} from '@mantine/core';
import {UseFormReturnType} from '@mantine/form';
import {Client4} from 'mattermost-redux/client';
import React, {ComponentPropsWithoutRef, forwardRef, useEffect, useState} from 'react';

import {UserProfile} from 'mattermost-redux/types/users';

import {getTeammateNameDisplaySetting} from 'mattermost-redux/selectors/entities/preferences';

import {GlobalState} from 'mattermost-redux/types/store';

import {useSelector} from 'react-redux';

import {getUser} from 'mattermost-redux/selectors/entities/users';

import {Customer} from '@/types/customers';
import {useDebounceSearch} from '@/hooks/debounce';
import {returnFormattedUsers} from '@/helpers/users';
type SelectProfileItem = {
    user: ReturnType<typeof returnFormattedUsers>[number];
    value: string;
    label: string;
}
interface ItemProps extends ComponentPropsWithoutRef<'div'>, SelectProfileItem {}
const SelectItem = forwardRef<HTMLDivElement, ItemProps>(
    ({user, ...others}: ItemProps, ref) => {
        return (
            <div
                ref={ref}
                id={`profile-${user.id}`}
                {...others}
            >
                <Group noWrap={true}>
                    <Avatar
                        src={user.image || ''}
                        radius='xl'
                    />
                    <div>
                        <Text size='sm'>{user.displayName}</Text>

                    </div>
                </Group>
            </div>
        );
    },
);

type FormUserSelectorParams = {
    getInputProps: UseFormReturnType<Customer, (values: Customer) => Customer>['getInputProps'];
    label: string;
    placeholder?: string;
    formKey: keyof Customer;
}
const FormUserSelector = (
    {
        getInputProps,
        label,
        placeholder,
        formKey,
    }: FormUserSelectorParams,
) => {
    const teamnameNameDisplaySetting = useSelector<GlobalState, string | undefined>(getTeammateNameDisplaySetting) || '';

    const [data, setData] = useState<SelectProfileItem[]>([]);
    const [selectedUser, setSelectedUser] = useState<SelectProfileItem | null>(null);

    const searchForUsers = async (term: string) => {
        if (!term) {
            setData([]);
            return;
        }

        if (selectedUser && selectedUser.user.displayName === term) {
            return;
        }

        await Client4.searchUsers(term, {limit: 20}).
            then((users) => {
                if (!users) {
                    return;
                }

                const formatted = returnFormattedUsers(users, teamnameNameDisplaySetting);
                setData(
                    formatted.map((user) => ({
                        user,
                        value: user.id,
                        label: user.displayName,
                    })),
                );
            });
    };

    const user = useSelector<GlobalState, UserProfile>((state) => getUser(state, getInputProps(formKey).value));

    useEffect(() => {
        if (user) {
            const formatted = returnFormattedUsers([user], teamnameNameDisplaySetting).map((u) => ({
                user: u,
                value: u.id,
                label: u.displayName,
            }));

            setSelectedUser(formatted[0]);
            setData(formatted);
        }
    }, [user, teamnameNameDisplaySetting]);

    const [, setSearchTermState] = useDebounceSearch(searchForUsers);

    return (
        <Select
            label={label}
            placeholder={placeholder || label}
            searchable={true}
            nothingFound='No options'
            onSearchChange={setSearchTermState}
            data={data}
            itemComponent={SelectItem}
            maxDropdownHeight={600}
            clearable={true}
            styles={() => ({
                root: {
                    width: '100%',
                },
                label: {
                    fontSize: '14px',
                    fontWeight: 600,
                },
                dropdown: {
                    border: '1px solid rgba(var(--center-channel-color-rgb), 0.2)',
                    background: 'var(--center-channel-bg)',
                    boxShadow: '0 0 0 1px hsla(0,0%,0%,0.1),0 4px 11px hsla(0,0%,0%,0.1)',
                    borderRadius: '4px',

                    // marginTop: '8px',
                    // marginBottom: '8px',
                },
                input: {
                    height: '40px',

                    // paddingLeft: '50px',
                    color: 'var(--center-channel-color)',
                    background: 'var(--center-channel-bg)',
                    border: '1px solid #ccc',

                    // padding: '6px 12px',
                    borderColor: 'rgba(var(--center-channel-color-rgb), 0.16)',
                    borderRadius: '4px',
                    lineHeight: '1.42857143',
                    transition: 'border-color ease-in-out .15s, box-shadow ease-in-out .15s, -webkit-box-shadow ease-in-out .15s',
                    '&:focus': {
                        borderColor: 'rgba(var(--button-bg-rgb), 1)',
                        boxShadow: '0 0 0 1px var(--button-bg)',
                    },
                },
                item: {
                    background: 'transparent',
                    height: '40px',

                    // applies styles to hovered item (with mouse or keyboard)
                    '&[data-hovered]': {
                        background: 'rgba(var(--button-bg-rgb), 0.16)',
                    },

                    // applies styles to selected item
                    '&[data-selected]': {
                        background: 'var(--button-bg)',
                        color: 'var(--button-color)',
                        '&, &:hover': {
                            background: 'rgba(var(--button-bg-rgb), 0.88)',
                        },
                    },
                },
                rightSection: {
                    color: 'var(--center-channel-color)',
                    width: '40px',
                },

            })}
            filter={() => true}
            {...getInputProps(formKey)}
        />
    );
};

export {
    FormUserSelector,
};
