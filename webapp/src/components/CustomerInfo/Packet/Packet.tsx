import React, {useState} from 'react';

import {Button, Divider, Title} from '@mantine/core';

import {PencilOutlineIcon} from '@mattermost/compass-icons/components';

import {useForm} from '@mantine/form';

import {CustomerPacketValues} from '@/types/customers';
import {Container} from '@/components/Group';

import {Item} from './Components/Item';
import {LicensedTo, MobileApp, PacketProvider, ProductsInUse, defaultPacket, transformKeys} from './Components/PacketItems';

type Params = {
    packet: CustomerPacketValues | undefined;
}
const CustomerInfoPacket = ({
    packet = defaultPacket,
}: Params) => {
    const [editing, setEditing] = useState(false);

    const {getInputProps} = useForm<CustomerPacketValues>({
        initialValues: packet,
    });

    return (
        <Container>
            <PacketProvider
                packet={packet}
                editing={editing}
                getInputProps={getInputProps}
            >
                {!editing && (
                    <Button
                        variant='subtle'
                        className='color--link'
                        leftIcon={<PencilOutlineIcon/>}
                        color='var(--center-channel-color)'
                        onClick={() => setEditing(true)}
                        sx={{
                            position: 'absolute',
                            top: '0',
                            right: '0',
                            margin: '1em',
                            '&:hover': {
                                backgroundColor: 'transparent',
                                textDecoration: 'underline',
                            },
                        }}
                    >
                        {'Edit'}
                    </Button>
                )}
                <div
                    style={{
                        display: 'flex',
                        flexDirection: 'column',
                        gap: '.5em',
                        padding: '4px 0',
                        width: '100%',
                        alignContent: 'flex-start',
                    }}
                >
                    <Title order={2}>{'General'}</Title>
                    <LicensedTo/>
                    <MobileApp/>
                    <ProductsInUse/>
                </div>
                <Divider
                    my='sm'
                    style={{width: '100%'}}
                />
                {/* // make this so it's broken into logical types of information. */}
                <div
                    style={{
                        display: 'flex',
                        flexDirection: 'column',
                        gap: '.5em',
                        padding: '4px 0',
                        width: '100%',
                        alignContent: 'flex-start',
                    }}
                >
                    <Title order={2}>{'Metrics'}</Title>
                    <Item
                        title={transformKeys.activeUsers}
                        value={packet.activeUsers.toLocaleString()}
                    />
                    <Item
                        title={transformKeys.dailyActiveUsers}
                        value={packet.dailyActiveUsers.toLocaleString()}
                    />
                    <Item
                        title={transformKeys.monthlyActiveUsers}
                        value={packet.monthlyActiveUsers.toLocaleString()}
                    />
                    <Item
                        title={transformKeys.totalChannels}
                        value={packet.totalChannels.toLocaleString()}
                    />
                    <Item
                        title={transformKeys.totalPosts}
                        value={packet.totalPosts.toLocaleString()}
                    />
                    <Item
                        title={transformKeys.totalTeams}
                        value={packet.totalTeams.toLocaleString()}
                    />
                    <Item
                        title={transformKeys.inactiveUserCount}
                        value={packet.inactiveUserCount.toLocaleString()}
                    />
                    <Item
                        title={transformKeys.licenseSupportedUsers}
                        value={packet.licenseSupportedUsers.toLocaleString()}
                    />
                </div>
                <Divider
                    my='sm'
                    style={{width: '100%'}}
                />
                <div
                    style={{
                        display: 'flex',
                        flexDirection: 'column',
                        gap: '.5em',
                        padding: '4px 0',
                        width: '100%',
                        alignContent: 'flex-start',
                    }}
                >
                    <Title order={2}>{'Environment'}</Title>
                    <Item
                        title={transformKeys.databaseType}
                        value={packet.databaseType}
                    />
                    <Item
                        title={transformKeys.serverOS}
                        value={packet.serverOS}
                    />
                    <Item
                        title={transformKeys.serverArch}
                        value={packet.serverArch}
                    />
                    <Item
                        title={transformKeys.version}
                        value={packet.version}
                    />
                    <Item
                        title={transformKeys.databaseSchemaVersion}
                        value={packet.databaseSchemaVersion}
                    />
                    <Item
                        title={transformKeys.databaseType}
                        value={packet.databaseType}
                    />
                    <Item
                        title={transformKeys.databaseVersion}
                        value={packet.databaseVersion}
                    />
                    <Item
                        title={transformKeys.elasticServerVersion}
                        value={packet.elasticServerVersion}
                    />
                    <Item
                        title={transformKeys.fileDriver}
                        value={packet.fileDriver}
                    />
                    <Item
                        title={transformKeys.hostingType}
                        value={packet.hostingType}
                    />
                    <Item
                        title={transformKeys.metrics}
                        value={(packet.metrics) ? 'true' : 'false'}
                    />
                    <Item
                        title={transformKeys.metricService}
                        value={packet.metricService}
                    />
                    <Item
                        title={transformKeys.ldapProvider}
                        value={packet.ldapProvider}
                    />
                    <Item
                        title={transformKeys.samlProvider}
                        value={packet.samlProvider}
                    />
                </div>
            </PacketProvider>
        </Container>
    );
};

export {
    CustomerInfoPacket,
};
