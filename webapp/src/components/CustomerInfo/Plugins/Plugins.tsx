import React from 'react';

import {Code} from '@mantine/core';

import {CustomerPluginValues} from '@/types/customers';
import {CenteredText} from '@/components/CenteredText';
import {PageHeader} from '../PageHeader';

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
            <PageHeader text='Plugins'/>
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
