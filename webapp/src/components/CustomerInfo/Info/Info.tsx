import React from 'react';

import {Code} from '@mantine/core';

import {Customer} from '@/types/customers';
import {CenteredText} from '@/components/CenteredText';

type Params = {
    customer: Customer | null;
}
const CustomerInfoProfile = ({
    customer,
}: Params) => {
    if (!customer) {
        return (
            <CenteredText
                message={'No customer information found.'}
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
            >{JSON.stringify(customer, null, 4)}</Code>
        </>

    );
};

export {
    CustomerInfoProfile,
};