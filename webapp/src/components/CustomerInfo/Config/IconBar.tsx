import {ActionIcon, Tooltip} from '@mantine/core';
import {CancelIcon, ContentCopyIcon, PencilOutlineIcon} from '@mattermost/compass-icons/components';
import React from 'react';

import {SaveIcon} from '@/components/SaveIcon';

type Params = {
    startEditing: () => void;
    cancelEditing: () => void;
    saveEdits: () => void;
    editing: boolean;
    copyToClipboard: () => void;
};

const IconBar = (
    {
        startEditing,
        cancelEditing,
        saveEdits,
        editing,
        copyToClipboard,
    }: Params,
) => {
    return (
        <div
            style={{
                position: 'absolute',
                right: '15px',
                top: '10px',
                zIndex: 100,
                display: 'flex',
                flexDirection: 'row',
                alignItems: 'center',
                justifyContent: 'center',
                gap: '1em',
            }}
        >
            <Tooltip
                label='Edit'
                position='bottom'
            >
                <ActionIcon
                    onClick={startEditing}
                    sx={{
                        '&:hover': {
                            backgroundColor: 'rgba(var(--center-channel-color-rgb), 0.08)',
                            color: 'rgba(var(--center-channel-color-rgb), 0.72)',
                        },
                        height: '28px',
                        width: '28px',
                        padding: '4px',

                        // we don't need to show the editing icon when editing.
                        display: editing ? 'none' : '',
                    }}
                >
                    <PencilOutlineIcon size='1.5em'/>
                </ActionIcon>
            </Tooltip>
            <Tooltip
                label='Cancel'
                position='bottom'
            >
                <ActionIcon
                    onClick={cancelEditing}
                    sx={{
                        '&:hover': {
                            backgroundColor: 'rgba(var(--center-channel-color-rgb), 0.08)',
                            color: 'rgba(var(--center-channel-color-rgb), 0.72)',
                        },
                        height: '28px',
                        width: '28px',
                        padding: '4px',

                        // we don't need to show the cancel icon when editing.
                        display: editing ? '' : 'none',
                    }}
                >
                    <CancelIcon size='1.5em'/>
                </ActionIcon>
            </Tooltip>
            <Tooltip
                label='Save'
                position='bottom'
            >
                <ActionIcon
                    onClick={saveEdits}
                    sx={{
                        '&:hover': {
                            backgroundColor: 'rgba(var(--center-channel-color-rgb), 0.08)',
                            color: 'rgba(var(--center-channel-color-rgb), 0.72)',
                        },
                        height: '28px',
                        width: '28px',
                        padding: '4px',

                        // we don't need to show the editing icon when editing.
                        display: editing ? '' : 'none',
                    }}
                >
                    <SaveIcon size='1.5em'/>
                </ActionIcon>
            </Tooltip>
            <Tooltip
                label='Copy to clipboard'
                position='bottom'
            >
                <ActionIcon
                    onClick={copyToClipboard}
                    sx={{
                        '&:hover': {
                            backgroundColor: 'rgba(var(--center-channel-color-rgb), 0.08)',
                            color: 'rgba(var(--center-channel-color-rgb), 0.72)',
                        },
                        height: '28px',
                        width: '28px',
                        padding: '4px',
                    }}
                >
                    <ContentCopyIcon size='1.5em'/>
                </ActionIcon>
            </Tooltip>
        </div>
    );
};

export {
    IconBar,
};
