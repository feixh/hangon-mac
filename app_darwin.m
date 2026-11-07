#import <Cocoa/Cocoa.h>

#include "app.h"
#include "_cgo_export.h"

void requestStop(void) {
    dispatch_async(dispatch_get_main_queue(), ^{
        [NSApp abortModal];
        [NSApp stop:nil];
        // Wake the event loop so run returns without another user event.
        [NSApp postEvent:[NSEvent otherEventWithType:NSEventTypeApplicationDefined
                                          location:NSZeroPoint modifierFlags:0
                                         timestamp:0 windowNumber:0 context:nil
                                           subtype:0 data1:0 data2:0]
                 atStart:YES];
    });
}

@interface HangOnDelegate : NSObject <NSApplicationDelegate>
@property(strong) NSStatusItem *statusItem;
@property(assign) AppConfig config;
@property(strong) NSData *icon;
@end

@implementation HangOnDelegate

- (void)applicationDidFinishLaunching:(NSNotification *)notification {
    (void)notification;
    self.statusItem = [[NSStatusBar systemStatusBar]
        statusItemWithLength:NSSquareStatusItemLength];
    NSImage *image = [[NSImage alloc] initWithData:self.icon];
    image.size = NSMakeSize(22, 22);
    image.template = YES;
    self.statusItem.button.image = image;
    self.statusItem.button.toolTip = @(self.config.tooltip);
    self.statusItem.button.accessibilityLabel = @(self.config.accessibility);

    NSMenu *menu = [[NSMenu alloc] initWithTitle:@(self.config.title)];
    NSMenuItem *status = [[NSMenuItem alloc] initWithTitle:@(self.config.status)
                                                 action:nil keyEquivalent:@""];
    status.enabled = NO;
    [menu addItem:status];
    [menu addItem:[NSMenuItem separatorItem]];
    NSMenuItem *stop = [[NSMenuItem alloc] initWithTitle:@(self.config.stop)
                                               action:@selector(stop:) keyEquivalent:@"q"];
    stop.target = self;
    [menu addItem:stop];
    self.statusItem.menu = menu;
    [[[NSWorkspace sharedWorkspace] notificationCenter]
        addObserver:self selector:@selector(machineDidWake:)
               name:NSWorkspaceDidWakeNotification object:nil];
}

- (void)stop:(id)sender {
    (void)sender;
    goStop();
}

- (void)machineDidWake:(NSNotification *)notification {
    (void)notification;
    dispatch_async(dispatch_get_main_queue(), ^{
        goDidWake();
    });
}

@end

int showPrompt(const char *title, const char *detail, const char *keep, const char *stop) {
    @autoreleasepool {
        [NSApplication sharedApplication];
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
        NSAlert *alert = [[NSAlert alloc] init];
        alert.messageText = @(title);
        alert.informativeText = @(detail);
        alert.alertStyle = NSAlertStyleInformational;
        [alert addButtonWithTitle:@(keep)];
        if (stop) {
            [alert addButtonWithTitle:@(stop)];
        }
        [NSApp activateIgnoringOtherApps:YES];
        return [alert runModal] == NSAlertSecondButtonReturn;
    }
}

void runApp(AppConfig config, const void *icon, size_t iconLength) {
    @autoreleasepool {
        [NSApplication sharedApplication];
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
        HangOnDelegate *delegate = [[HangOnDelegate alloc] init];
        delegate.config = config;
        delegate.icon = [NSData dataWithBytes:icon length:iconLength];
        NSApp.delegate = delegate;
        [NSApp run];
        [[[NSWorkspace sharedWorkspace] notificationCenter] removeObserver:delegate];
        [[NSStatusBar systemStatusBar] removeStatusItem:delegate.statusItem];
        NSApp.delegate = nil;
    }
}
