import React from 'react';

import {Code} from '@mantine/core';

import {CustomerConfigValues} from '@/types/customers';
import {CenteredText} from '@/components/CenteredText';

type Params = {
    config: CustomerConfigValues | null;
}
const CustomerInfoConfig = ({
    config,
}: Params) => {
    if (!config) {
        return (
            <CenteredText
                message={'No config information found.'}
            />
        );
    }
    return (
        <>
            <Code
                block={true}
                style={{
                    width: '100%',
                }}
            >{JSON.stringify(config, null, 4)}</Code>
        </>

    );
};

export {
    CustomerInfoConfig,
};
