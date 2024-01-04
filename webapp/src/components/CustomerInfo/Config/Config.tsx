import React, {useCallback, useEffect, useRef, useState} from 'react';

import {Prism} from '@mantine/prism';

import {createStyles} from '@mantine/styles';

import {Textarea} from '@mantine/core';

import {CustomerConfigValues} from '@/types/customers';
import {Container} from '@/components/Group';

import {ConfigSearchBar} from './SearchBar';
import {IconBar} from './IconBar';

const useStylesPrism = createStyles(() => ({
    root: {
        border: '1px solid rgba(var(--center-channel-color-rgb), 0.2)',
        width: '100%',
        overflow: 'scroll',
    },
    code: {
        color: 'var(--center-channel-color)',
        background: 'var(--center-channel-bg)',
        fontSize: '13px',
        border: 'none',
    },
    copy: {
        display: 'none',
    },
}));

const useStylesEdit = createStyles(() => ({
    root: {
        width: '100%',
        overflow: 'scroll',
        height: '100%',
    },
    input: {
        color: 'var(--center-channel-color)',
        background: 'var(--center-channel-bg)',
        fontSize: '13px',
        border: 'none',
        height: '100%',
    },
    wrapper: {
        height: '100%',
    },
}));

  type Params = {
      config: CustomerConfigValues | null;
      save: (config: CustomerConfigValues) => void;
  }

type LineHighlight = { [key: number]: { color: string} };

const stringifyConfig = (config: CustomerConfigValues | null | string) => {
    return JSON.stringify(config || {}, null, 4);
};
const CustomerInfoConfig = ({
    config,
    save,
}: Params) => {
    const {classes: prismClasses} = useStylesPrism();
    const {classes: editClasses} = useStylesEdit();
    const prismRef = useRef<HTMLDivElement>(null); // create a ref
    const parsedConfig = stringifyConfig(config);

    const [searchValue, setSearchValue] = useState('');
    const [lineHighlights, setLineHighlights] = useState<LineHighlight>({});
    const [editing, setEditing] = useState(false);
    const [editingState, setEditingState] = useState(''); // this is the value of the text area
    const [error, setError] = useState('');

    useEffect(() => {
        setEditingState(parsedConfig);
    }, [parsedConfig]);

    const scrollLineIntoView = useCallback((lineNumber: number) => {
        if (!prismRef.current) {
            return;
        }
        const lineElement = prismRef.current.querySelector(`.token-line:nth-child(${lineNumber})`);
        if (lineElement) {
            lineElement.scrollIntoView({behavior: 'auto', block: 'center'});
        }
    }, [prismRef]);

    const updateHighlightAndScroll = useCallback((conf: typeof config, term: string) => {
        // no config, search term, or prism ref, so we can't do anything.
        if (!conf || !term) {
            setLineHighlights({});
            scrollLineIntoView(0);
            return;
        }

        // splitting into array of config values to search through easier.
        const lines = stringifyConfig(conf).split('\n');
        const indexes: number[] = [];
        lines.forEach((line, index) => {
            if (line.toLowerCase().includes(term.toLowerCase())) {
                indexes.push(index);
            }
        });
        if (indexes.length === 0) {
            setLineHighlights({});
            scrollLineIntoView(0);
            return;
        }

        const highlightLines: LineHighlight = {};
        for (const num of indexes) {
            highlightLines[num + 1] = {color: 'green'};
        }

        setLineHighlights(highlightLines);
        scrollLineIntoView(indexes[0] + 1);
    }, [setLineHighlights, scrollLineIntoView]);

    useEffect(() => {
        updateHighlightAndScroll(config, searchValue);
    }, [searchValue, config, updateHighlightAndScroll]);

    const handleSearchChange = (value: string) => {
        setSearchValue(value);
    };

    const copyToClipboard = () => {
        if (editing) {
            navigator.clipboard.writeText(stringifyConfig(editingState));
            return;
        }
        navigator.clipboard.writeText(parsedConfig);
    };

    const onEditChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        setEditingState(e.target.value);
        if (error) {
            setError('');
        }
    };

    const startEditing = () => {
        setEditing(true);
    };

    const cancelEditing = () => {
        setEditing(false);
        setEditingState(parsedConfig);
    };

    const isValidJson = (jsonString: string): boolean => {
        try {
            JSON.parse(jsonString);
        } catch (e) {
            return false;
        }
        return true;
    };

    const saveEdits = () => {
        if (editingState !== parsedConfig) {
            if (!isValidJson(editingState)) {
                setError('Invalid JSON');
                return;
            }
            save(JSON.parse(editingState));
        }
        setEditing(false);
    };

    const textareaRef = useRef<HTMLTextAreaElement>(null);

    const searchForStringInEdit = useCallback((term: string) => {
        const textArea = textareaRef.current;
        if (!textArea) {
            return;
        }
        const fullText = textArea.value;

        const specificString = term;
        const index = fullText.toLowerCase().indexOf(specificString.toLowerCase());

        if (index !== -1) {
            const activeElement = document.activeElement;
            const lineHeight = 1.5 * 13; // replace with your actual line height
            const lines = fullText.substring(0, index).split('\n');
            const lineNumber = lines.length;
            const scrollTop = lineHeight * lineNumber;
            textArea.scrollTop = scrollTop;
            textArea.setSelectionRange(index, index + specificString.length);

            if (activeElement instanceof HTMLElement) {
                activeElement.focus();
            }
        }
    }, []);

    useEffect(() => {
        if (searchValue && editing) {
            searchForStringInEdit(searchValue);
        }
    }, [searchForStringInEdit, searchValue, editing]);

    return (
        <Container>
            <ConfigSearchBar
                handleSearchChange={handleSearchChange}
                disabled={false}
            />

            <div
                style={{
                    height: '100%',
                    width: '100%',
                    overflow: 'hidden',
                    display: 'flex',
                    flexDirection: 'column',
                    position: 'relative',
                }}
            >
                <IconBar
                    startEditing={startEditing}
                    cancelEditing={cancelEditing}
                    saveEdits={saveEdits}
                    editing={editing}
                    copyToClipboard={copyToClipboard}
                />
                {
                    editing ? (
                        <Textarea
                            ref={textareaRef}
                            bg={'transparent'}
                            classNames={editClasses}
                            value={editingState}
                            onChange={onEditChange}
                            onFocus={() => searchForStringInEdit(searchValue)}
                            styles={{
                                root: {
                                    border: error ? '1px solid var(--error-text)' : '1px solid rgba(var(--center-channel-color-rgb), 0.2)',
                                },
                            }}
                        />
                    ) : (
                        <Prism
                            ref={prismRef}
                            withLineNumbers={true}
                            highlightLines={lineHighlights}
                            language='json'
                            contentEditable={false}
                            classNames={prismClasses}
                        >
                            {parsedConfig}
                        </Prism>
                    )
                }

            </div>
        </Container>

    );
};

export {
    CustomerInfoConfig,
};
