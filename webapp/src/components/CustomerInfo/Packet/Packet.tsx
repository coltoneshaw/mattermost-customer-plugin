import React, {useEffect, useState} from 'react';

import {Button, Divider, Group, Title} from '@mantine/core';

import {PencilOutlineIcon} from '@mattermost/compass-icons/components';

import {useForm} from '@mantine/form';

import {CustomerPacketValues} from '@/types/customers';
import {Container} from '@/components/Group';

import {PacketNumberField, PacketProvider, PacketSelect, PacketTextInput, ProductsInUse, defaultPacket} from './Components/PacketItems';

type Params = {
    packet: CustomerPacketValues | undefined;
    save: (values: CustomerPacketValues) => void;
}
const CustomerInfoPacket = ({
    packet = defaultPacket,
    save,
}: Params) => {
    const [editing, setEditing] = useState(false);

    const {getInputProps, setValues, resetDirty, resetTouched, values} = useForm<CustomerPacketValues>({
        initialValues: defaultPacket,
        validateInputOnBlur: true,
        clearInputErrorOnChange: true,
        validate: {
            version: (value) => {
                const versionRegex = /^\d+\.\d+\.\d+$/;
                if (!versionRegex.test(value)) {
                    return 'Version must be in x.x.x format';
                }
                return true;
            },

            elasticServerVersion: (value) => {
                const versionRegex = /^\d+\.\d+\.\d+$/;
                if (!versionRegex.test(value)) {
                    return 'Version must be in x.x.x format';
                }
                return true;
            },
        },
    });

    // resetting the form is anything new comes in.
    useEffect(() => {
        if (packet) {
            setValues(packet);
            resetDirty(packet);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [packet]);

    const saveValues = () => {
        save(values);
        setEditing(false);
    };
    const resetForm = () => {
        setValues(packet);
        resetDirty(packet);
        resetTouched();
        setEditing(false);
    };

    return (
        <>
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
                        <PacketTextInput packetKey='licensedTo'/>
                        <PacketSelect
                            packetKey='mobileApp'
                            data={[
                                {value: 'using', label: 'Using'},
                                {value: 'not using', label: 'Not using'},
                                {value: 'not-using security', label: 'Not using security issues'},
                                {value: 'unknown', label: 'unknown'},
                            ]}
                        />
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
                        <PacketNumberField packetKey={'activeUsers'}/>
                        <PacketNumberField packetKey={'dailyActiveUsers'}/>
                        <PacketNumberField packetKey={'monthlyActiveUsers'}/>
                        <PacketNumberField packetKey={'inactiveUserCount'}/>
                        <PacketNumberField packetKey={'licenseSupportedUsers'}/>
                        <PacketNumberField packetKey={'totalPosts'}/>
                        <PacketNumberField packetKey={'totalChannels'}/>
                        <PacketNumberField packetKey={'totalTeams'}/>
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
                        <PacketSelect
                            packetKey={'deploymentType'}
                            data={[
                                {value: 'aks', label: 'AKS'},
                                {value: 'kubernetes', label: 'Kubernetes'},
                                {value: 'docker', label: 'Docker'},
                                {value: 'tarball', label: 'Tarball'},
                                {value: 'mattermost-cloud', label: 'Mattermost Cloud'},
                                {value: 'vm', label: 'VM'},
                                {value: 'other', label: 'Other'},
                            ]}
                        />
                        <PacketSelect
                            packetKey={'hostingType'}
                            data={[
                                {value: 'aws', label: 'AWS'},
                                {value: 'aws-gov', label: 'AWS Gov Cloud'},
                                {value: 'azure', label: 'Azure'},
                                {value: 'gcp', label: 'Google Cloud'},
                                {value: 'digital-ocean', label: 'Digial Ocean'},
                                {value: 'mattermost-cloud', label: 'Mattermost Cloud'},
                                {value: 'on-prem', label: 'On Prem'},
                            ]}
                        />
                        <PacketSelect
                            packetKey={'serverOS'}
                            data={[
                                {value: 'linux', label: 'Linux'},
                                {value: 'darwin', label: 'Darwin'},
                                {value: 'windows', label: 'Windows'},
                            ]}
                        />
                        <PacketTextInput
                            packetKey={'serverArch'}
                        />
                        <PacketTextInput
                            packetKey={'version'}
                        />
                        <PacketSelect
                            packetKey={'databaseType'}
                            data={[
                                {value: 'postgres', label: 'PostgreSQL'},
                                {value: 'mysql', label: 'MySQL'},
                                {value: 'mariadb', label: 'MariaDB'},
                                {value: 'other', label: 'Other'},
                            ]}
                        />
                        <PacketTextInput
                            packetKey='databaseSchemaVersion'
                        />
                        <PacketTextInput
                            packetKey={'databaseVersion'}
                        />
                        <PacketTextInput
                            packetKey={'elasticServerVersion'}
                        />
                        <PacketSelect
                            packetKey={'fileDriver'}
                            data={[
                                {value: 'amazons3', label: 'Amazon S3'},
                                {value: 'local', label: 'Local'},
                            ]}
                        />
                        <PacketTextInput
                            packetKey={'metricService'}
                        />
                        <PacketTextInput
                            packetKey={'ldapProvider'}
                        />
                        <PacketTextInput
                            packetKey={'samlProvider'}
                        />
                    </div>

                </PacketProvider>
            </Container>
            {
                editing && (
                    <Group
                        style={{
                            borderTop: '1px solid rgba(var(--center-channel-color-rgb), 0.08)',
                            padding: '1em',
                        }}
                    >
                        <button
                            className='btn btn-primary'
                            onClick={saveValues}
                        >
                            {'Save'}
                        </button>
                        <button
                            className='btn btn-tertiary'
                            onClick={resetForm}
                        >
                            {'Cancel'}
                        </button>
                    </Group>
                )
            }
        </>
    );
};

export {
    CustomerInfoPacket,
};
