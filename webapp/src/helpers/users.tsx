import {Client4} from 'mattermost-redux/client';
import {UserProfile} from 'mattermost-redux/types/users';
import {displayUsername} from 'mattermost-redux/utils/user_utils';

export const returnFormattedUsers = (users: UserProfile[], teammateNameDisplaySetting: string) => {
    return users.map((user) => {
        const displayName = displayUsername(user, teammateNameDisplaySetting);
        return {
            ...user,
            displayName,
            image: Client4.getProfilePictureUrl(user.id, user.last_picture_update),
        };
    });
};
