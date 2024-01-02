import {DateTime} from 'luxon';

export const getTimestamp = (epoch: number) => {
    return DateTime.fromMillis(epoch).toRelative();
};

