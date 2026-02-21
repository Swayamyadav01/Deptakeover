# Security Policy

## 🔒 Security Overview

DepTakeover is a security tool designed to identify supply chain vulnerabilities. We take the security of the tool itself seriously and appreciate the security research community's help in keeping it secure.

## 🚨 Supported Versions

We provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | ✅ Yes             |
| < 1.0   | ❌ No              |

## 🐛 Reporting a Vulnerability

### For DepTakeover Security Issues

If you discover a security vulnerability in DepTakeover itself, please help us maintain security by following responsible disclosure:

**⚠️ DO NOT create public GitHub issues for security vulnerabilities.**

Instead:

1. **Email us**: Send details to `security@deptakeover.dev` (or your contact email)
2. **Include**:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fixes (if any)
   - Your contact information

### Response Timeline

- **Initial Response**: Within 24-48 hours
- **Assessment**: Within 1 week  
- **Fix Development**: 1-4 weeks depending on severity
- **Public Disclosure**: After fix is released and users have time to update

### Severity Levels

**Critical (4-6 hours response)**
- Remote code execution
- Authentication bypass
- Data exfiltration capabilities

**High (24-48 hours response)**  
- Privilege escalation
- Information disclosure of sensitive data
- Local file inclusion/traversal

**Medium (1 week response)**
- Denial of service
- Information disclosure (non-sensitive)
- Logic flaws

**Low (2 weeks response)**
- Enhancement suggestions
- Minor information disclosure

## 🏆 Security Researcher Recognition

We believe in recognizing security researchers who help improve DepTakeover's security:

### Hall of Fame
Contributors who responsibly disclose security issues will be listed here (with permission):

*No security issues reported yet - be the first!*

### Rewards
While we don't offer monetary rewards, we provide:
- Public recognition (if desired)
- Mention in release notes
- Direct communication with maintainers
- Priority support for future issues

## 🛡️ Security Best Practices for Users

### When Using DepTakeover

1. **Keep Updated**: Always use the latest version
2. **Verify Downloads**: Check integrity of binaries from official releases
3. **Limited Scope**: Only scan repositories you have permission to analyze
4. **Secure Environment**: Run in isolated environments when possible
5. **Rate Limiting**: Tool includes built-in rate limiting, but be respectful of APIs

### For Bug Bounty Hunters

1. **Authorized Testing Only**: Ensure you have permission to scan targets
2. **Responsible Disclosure**: Follow target's vulnerability disclosure policy
3. **Verify Results**: Manually confirm package takeover possibilities
4. **Document Impact**: Clearly explain potential impact of findings
5. **Respect Scope**: Stay within defined bug bounty program scope

## 🔍 Security Features

DepTakeover includes several security features:

### Input Validation
- Repository URL sanitization
- Package name validation
- File path traversal protection

### Network Security
- HTTPS-only API communication
- Certificate validation
- Request timeout handling
- Rate limiting compliance

### Local Security
- No persistent storage of sensitive data
- Temporary file cleanup
- Limited file system access
- No elevated privileges required

## 🚫 Out of Scope

The following are generally NOT considered security vulnerabilities:

### Expected Behavior
- Rate limiting by external APIs (npm, PyPI, GitHub)
- Network connectivity issues
- Permission errors for private repositories
- Resource exhaustion on extremely large organizations

### User Error
- Running tool against unauthorized targets
- Misuse of scan results
- Failure to follow responsible disclosure

### Third-Party Issues
- Vulnerabilities in Go standard library
- GitHub API security issues  
- Package registry security issues
- Operating system vulnerabilities

## 📋 Security Checklist for Contributors

Before submitting code:

- [ ] Input validation for all user inputs
- [ ] Error handling that doesn't leak sensitive information
- [ ] No hardcoded credentials or tokens
- [ ] Secure HTTP client configuration
- [ ] Proper file handle cleanup
- [ ] Path traversal protection
- [ ] Rate limiting respect

## 📞 Additional Security Resources

- [Go Security Best Practices](https://golang.org/doc/security/)
- [OWASP Secure Coding Practices](https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/)
- [GitHub Security Advisories](https://docs.github.com/en/code-security/security-advisories)

---

Thank you for helping keep DepTakeover and its users secure! 🛡️