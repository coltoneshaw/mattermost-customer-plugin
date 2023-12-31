import React from 'react';

import {Code} from '@mantine/core';

import {CustomerPacketValues} from '@/types/customers';
import {CenteredText} from '@/components/CenteredText';
import {Container} from '@/components/Group';

type Params = {
    packet: CustomerPacketValues | null;
}
const CustomerInfoPacket = ({
    packet,
}: Params) => {
    if (!packet) {
        return (
            <CenteredText
                message={'No packet information found.'}
            />
        );
    }
    return (
        <Container>
            <Code
                block={true}
                style={{
                    width: '100%',
                }}
            >{JSON.stringify(packet, null, 4)}</Code>
        </Container>
    );
};

export {
    CustomerInfoPacket,
};
