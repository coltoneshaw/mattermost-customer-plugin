import React from 'react';

import {Code} from '@mantine/core';

import {CustomerPluginValues} from '@/types/customers';
import {CenteredText} from '@/components/CenteredText';

type Params = {
    plugins: CustomerPluginValues[] | null;
}
const CustomerInfoPlugins = ({
    plugins,
}: Params) => {
    if (!plugins || plugins.length === 0) {
        return (
            <CenteredText
                message={'No plugin information found.'}
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
            >{JSON.stringify(plugins, null, 4)}</Code>
        </>
    );
};

export {
    CustomerInfoPlugins,
};
