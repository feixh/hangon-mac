#ifndef HANGON_APP_H
#define HANGON_APP_H

#include <stddef.h>

typedef struct {
    const char *title;
    const char *status;
    const char *stop;
    const char *tooltip;
    const char *accessibility;
} AppConfig;

void runApp(AppConfig config, const void *icon, size_t iconLength);
int showPrompt(const char *title, const char *detail, const char *keep, const char *stop);
void requestStop(void);

#endif
