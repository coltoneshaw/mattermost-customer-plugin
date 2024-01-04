import {Avatar, Tooltip} from '@mantine/core';
import React from 'react';

import {useSelector} from 'react-redux';

// import {GlobalState} from '@mattermost/types/lib/store';
import {getUser} from 'mattermost-redux/selectors/entities/users';
import {GlobalState} from 'mattermost-redux/types/store';
import {getTeammateNameDisplaySetting} from 'mattermost-redux/selectors/entities/preferences';

import {UserProfile} from 'mattermost-redux/types/users';

import {Customer} from '@/types/customers';
import {returnFormattedUsers} from '@/helpers/users';

type ProfileParams = {
    userId: string;
    role: string;
}
const Profile = ({
    userId,
    role,
}: ProfileParams) => {
    const teamnameNameDisplaySetting = useSelector<GlobalState, string | undefined>(getTeammateNameDisplaySetting) || '';

    const user = useSelector<GlobalState, UserProfile>((state) => getUser(state, userId));

    const formattedUser = returnFormattedUsers([user], teamnameNameDisplaySetting)[0];

    return (
        <Tooltip
            withArrow={true}
            label={formattedUser.displayName + ' - ' + role}
        >
            <Avatar
                src={formattedUser.image || ''}
                radius='xl'
                sx={{
                    height: '32px',
                    width: '32px',
                    border: '2px dotted rgba(0, 0, 0, 0)',
                    backgroundColor: 'var(--center-channel-bg)',
                    '&:hover': {
                        zIndex: 100,
                    },
                }}
            />
        </Tooltip>

    );
};

type Roles = Pick<Customer, 'accountExecutive' | 'customerSuccessManager' | 'technicalAccountManager'>
type UserListParams = Roles
const UserList = ({
    accountExecutive,
    customerSuccessManager,
    technicalAccountManager,
}: UserListParams) => {
    return (
        <Avatar.Group spacing='sm'>
            {accountExecutive && (
                <Profile
                    userId={accountExecutive}
                    role='Account Executive'
                />
            )}
            {customerSuccessManager && (
                <Profile
                    userId={customerSuccessManager}
                    role='Customer Success Manager'
                />
            )}
            {technicalAccountManager && (
                <Profile
                    userId={technicalAccountManager}
                    role='Technical Account Manager'

                />
            )}
        </Avatar.Group>
    );
};

export {
    UserList,
};
