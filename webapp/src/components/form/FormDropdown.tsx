import React from 'react';

import {Select} from '@mantine/core';

import {FormTextInputParams} from './types';
import {formSelectInputStyles} from './styles';

const FormDropdown = <T, >({
    getInputProps,
    label,
    placeholder,
    formKey,
    data,
}: FormTextInputParams<T> & { data: { value: string, label: string }[]},
) => {
    const {classes} = formSelectInputStyles();
    return (
        <Select
            bg={'transparent'}
            clearable={true}
            maxDropdownHeight={600}
            classNames={classes}
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
