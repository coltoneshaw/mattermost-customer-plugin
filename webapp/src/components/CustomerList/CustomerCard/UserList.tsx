import {Avatar, Tooltip} from '@mantine/core';
import React, {useEffect} from 'react';

import {useDispatch, useSelector} from 'react-redux';

// import {GlobalState} from '@mattermost/types/lib/store';
import {getUser} from 'mattermost-redux/selectors/entities/users';
import {GlobalState} from 'mattermost-redux/types/store';
import {getTeammateNameDisplaySetting} from 'mattermost-redux/selectors/entities/preferences';
import {displayUsername} from 'mattermost-redux/utils/user_utils';
import {getUser as fetchUser} from 'mattermost-redux/actions/users';
import {Client4} from 'mattermost-redux/client';
import {UserProfile} from 'mattermost-redux/types/users';

import {Customer} from '@/types/customers';

type ProfileParams = {
    userId: string;
    role: string;
}
const Profile = ({
    userId,
    role,
}: ProfileParams) => {
    const dispatch = useDispatch();

    const user = useSelector<GlobalState, UserProfile>((state) => getUser(state, userId));
    const teamnameNameDisplaySetting = useSelector<GlobalState, string | undefined>(getTeammateNameDisplaySetting) || '';
    useEffect(() => {
        if (!user || !user.id || !dispatch) {
            dispatch(fetchUser(userId));
        }
    }, [userId, dispatch, user]);

    let name = null;
    let profileUri = null;
    if (user) {
        const preferredName = displayUsername(user, teamnameNameDisplaySetting);
        name = preferredName;
        profileUri = Client4.getProfilePictureUrl(userId, user.last_picture_update);
    }

    return (
        <Tooltip
            withArrow={true}
            label={name + ' - ' + role}
        >
            <Avatar
                src={profileUri || ''}
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
