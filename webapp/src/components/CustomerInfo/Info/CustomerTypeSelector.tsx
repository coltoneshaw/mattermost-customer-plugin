import React from 'react';

import {Select} from '@mantine/core';

import {FormTextInputParams} from './Info';

const FormDropdown = ({
    getInputProps,
    label,
    placeholder,
    formKey,
    data,
}: FormTextInputParams & { data: { value: string, label: string }[]},
) => {
    return (
        <Select
            bg={'transparent'}
            clearable={true}
            maxDropdownHeight={600}
            styles={() => ({
                root: {
                    width: '100%',
                },
                label: {
                    fontSize: '14px',
                    fontWeight: 600,
                },
                dropdown: {
                    border: '1px solid rgba(var(--center-channel-color-rgb), 0.2)',
                    background: 'var(--center-channel-bg)',
                    boxShadow: '0 0 0 1px hsla(0,0%,0%,0.1),0 4px 11px hsla(0,0%,0%,0.1)',
                    borderRadius: '4px',
                    marginTop: '8px',
                    marginBottom: '8px',
                },
                input: {
                    height: '40px',
                    paddingLeft: '50px',
                    color: 'var(--center-channel-color)',
                    background: 'var(--center-channel-bg)',
                    border: '1px solid #ccc',
                    padding: '6px 12px',
                    borderColor: 'rgba(var(--center-channel-color-rgb), 0.16)',
                    borderRadius: '4px',
                    lineHeight: '1.42857143',
                    transition: 'border-color ease-in-out .15s, box-shadow ease-in-out .15s, -webkit-box-shadow ease-in-out .15s',
                    '&:focus': {
                        borderColor: 'rgba(var(--button-bg-rgb), 1)',
                        boxShadow: '0 0 0 1px var(--button-bg)',
                    },
                },
                item: {
                    background: 'transparent',
                    height: '40px',

                    // applies styles to hovered item (with mouse or keyboard)
                    '&[data-hovered]': {
                        background: 'rgba(var(--button-bg-rgb), 0.16)',
                    },

                    // applies styles to selected item
                    '&[data-selected]': {
                        background: 'var(--button-bg)',
                        color: 'var(--button-color)',
                        '&, &:hover': {
                            background: 'rgba(var(--button-bg-rgb), 0.88)',
                        },
                    },
                },
                rightSection: {
                    color: 'var(--center-channel-color)',
                    width: '40px',
                },
            })}
            placeholder={placeholder || label}
            label={label}
            data={data}
            {...getInputProps(formKey)}
        />
    );
};

export {
    FormDropdown,
};
